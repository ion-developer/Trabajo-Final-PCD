const jsonServer = require('json-server');
const server = jsonServer.create();
const router = jsonServer.router('db.json');
const middlewares = jsonServer.defaults();
const port = process.env.PORT || 3000;
const path = '/api/v1';
// server.setPath(path);
server.use(middlewares);
server.use(path,router);

server.listen(port );