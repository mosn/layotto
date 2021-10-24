const Koa = require('koa');
const app = new Koa();

// response
app.use(ctx => {
  console.log('%s %s, headers: %j', ctx.method, ctx.url, ctx.headers);
  ctx.body = 'Hello Koa';
});

app.listen(8889);
