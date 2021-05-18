package actuator

var singleton = New()

func GetDefault() *Actuator {
	return singleton
}
