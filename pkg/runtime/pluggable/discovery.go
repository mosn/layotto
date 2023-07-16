package pluggable

import (
	"context"
	"github.com/jhump/protoreflect/grpcreflect"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	reflectpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
	"io/fs"
	"mosn.io/layotto/components/configstores"
	"mosn.io/layotto/components/file"
	"mosn.io/layotto/components/hello"
	"mosn.io/layotto/components/oss"
	"mosn.io/layotto/components/rpc"
	mbindings "mosn.io/layotto/pkg/runtime/bindings"
	runtime_lock "mosn.io/layotto/pkg/runtime/lock"
	"mosn.io/layotto/pkg/runtime/pubsub"
	msecretstores "mosn.io/layotto/pkg/runtime/secretstores"
	runtime_sequencer "mosn.io/layotto/pkg/runtime/sequencer"
	"mosn.io/layotto/pkg/runtime/state"
	"mosn.io/pkg/log"
	"os"
	"path/filepath"
)

var onServiceDiscovered map[string]func(compType string, dialer GRPCConnectionDialer, c *Discovery)

func init() {
	onServiceDiscovered = make(map[string]func(compType string, dialer GRPCConnectionDialer, c *Discovery))
}

func AddServiceDiscoveryCallback(serviceName string, callbackFunc func(compType string, dialer GRPCConnectionDialer, c *Discovery)) {
	onServiceDiscovered[serviceName] = callbackFunc
}

const (
	defaultSocketFolder = "/tmp/runtime/component-sockets"
	SocketFolderEnvVar  = "LAYOTTO_COMPONENTS_SOCKETS_FOLDER"
)

type Discovery struct {
	Hellos        []*hello.HelloFactory
	ConfigStores  []*configstores.StoreFactory
	RPCs          []*rpc.Factory
	Files         []*file.FileFactory
	Oss           []*oss.Factory
	PubSubs       []*pubsub.Factory
	States        []*state.Factory
	Locks         []*runtime_lock.Factory
	Sequencers    []*runtime_sequencer.Factory
	OutputBinding []*mbindings.OutputBindingFactory
	InputBinding  []*mbindings.InputBindingFactory
	SecretStores  []*msecretstores.SecretStoresFactory

	// extension components

}

func Discover() (*Discovery, error) {
	d := new(Discovery)
	services, err := serviceDiscovery()
	if err != nil {
		return nil, err
	}

	// register services to discovery
	for _, s := range services {
		f, ok := onServiceDiscovered[s.protoRef]
		if !ok {
			continue
		}

		// use callback adding pluggable component to Discovery factories
		f(s.compType, s.dialer, d)
	}

	return d, nil
}

type service struct {
	protoRef string // the proto service name
	compType string // component type
	dialer   GRPCConnectionDialer
}

func GetSocketFolderPath() string {
	if v, ok := os.LookupEnv(SocketFolderEnvVar); ok {
		return v
	}
	return defaultSocketFolder
}

func serviceDiscovery() ([]service, error) {
	var services []service
	files, err := getComponentSocketFiles()
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		compType := removeFileExtension(f)
		dialer := socketDialer(f)
		grpcService, err := reflectGRPCServices(f)
		if err != nil {
			return nil, err
		}

		for _, s := range grpcService {
			services = append(services, service{
				protoRef: s,
				compType: compType,
				dialer:   dialer,
			})
		}
	}

	return services, nil
}

func getComponentSocketFiles() ([]string, error) {
	// the full filepath of the component socket file
	var socketFilepath []string

	// 1. get pluggable component socket file folder
	componentSocketPath := GetSocketFolderPath()

	// 2. check if the component socket path is existed or not.
	_, err := os.Stat(componentSocketPath)
	if ok := os.IsNotExist(err); ok {
		log.DefaultLogger.Infof("pluggable component socket folder [%s] not exist", componentSocketPath)
		return socketFilepath, nil
	}
	if err != nil {
		return nil, err
	}

	// 3. read socket files.
	socketFiles, err := os.ReadDir(componentSocketPath)
	if err != nil {
		return nil, err
	}

	for _, dirEntry := range socketFiles {
		// skip folder
		if ok := dirEntry.IsDir(); ok {
			continue
		}

		f, err := dirEntry.Info()
		if err != nil {
			return nil, err
		}

		// skip if file is not a socket type
		if ok := isSocketFileType(f); !ok {
			continue
		}

		fullPath := filepath.Join(componentSocketPath, dirEntry.Name())
		socketFilepath = append(socketFilepath, fullPath)

	}
	return socketFilepath, nil
}

func isSocketFileType(f os.FileInfo) bool {
	return f.Mode()&fs.ModeSocket != 0
}

func removeFileExtension(f string) string {
	// remove file path message
	filename := filepath.Base(f)
	// remove the file type extension message
	return filename[:len(filename)-len(filepath.Ext(filename))]
}

func reflectGRPCServices(socket string) ([]string, error) {
	// 1. connect to grpc service
	conn, err := socketDial(context.TODO(), socket)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// 2. reflect grpc connection
	client := grpcreflect.NewClientV1Alpha(context.TODO(), reflectpb.NewServerReflectionClient(conn))
	defer client.Reset()

	// 3. get service list
	services, err := client.ListServices()
	if err != nil {
		return nil, err
	}

	return services, nil
}

func socketDial(ctx context.Context, socket string, additionalOpts ...grpc.DialOption) (*grpc.ClientConn, error) {
	udsSocket := "unix://" + socket
	additionalOpts = append(additionalOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	grpcConn, err := grpc.DialContext(ctx, udsSocket, additionalOpts...)
	if err != nil {
		return nil, err
	}
	return grpcConn, nil
}

type GRPCConnectionDialer func(ctx context.Context, opts ...grpc.DialOption) (*grpc.ClientConn, error)

// create a socket-specific grpc connection via a closure
func socketDialer(socket string) GRPCConnectionDialer {
	return func(ctx context.Context, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
		return socketDial(ctx, socket, opts...)
	}
}
