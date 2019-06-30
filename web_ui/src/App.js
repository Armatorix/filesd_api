import React from 'react';
import { Button, Container, Row, Badge, Table, Navbar, Nav } from 'react-bootstrap'
import './App.css';
// TODO add multiple configs per config
var configs = [
  {
    id: "6bc661257fd81341d93fc741cf1999684baeacec",
    configs: [
      {
        targets: ["123.321.213.321:3213", "123.321.213.321:3213", "123.321.213.321:3213", "123.321.213.321:3213"],
        labels: { "1": "1lab", "2": "2lab" },
      },
      {
        targets: ["123.321.213.321:3213", "123.321.213.321:3213", "123.321.213.321:3213", "123.321.213.321:3213"],
        labels: { "1": "1lab", "2": "2lab" },
      }
    ]
  },
  {
    id: "2137",
    configs: [
      {
        targets: ["23.321.213.321:3213", "123.321.213.321:3213", "123.321.213.321:3213", "123.321.213.321:3213"],
        labels: { "1": "1lab", "2": "2lab" },
      },
      {
        targets: ["123.321.213.321:3213", "123.321.213.321:3213", "123.321.213.321:3213", "123.321.213.321:3213"],
        labels: { "1": "1lab", "2": "2lab" },
      },
      {
        targets: ["123.321.213.321:3213", "123.321.213.321:3213", "123.321.213.321:3213", "123.321.213.321:3213"],
        labels: { "1": "1lab", "2": "2lab" },
      }
    ]
  }
]


function buildSubTable(c) {
  console.log(c)
  var rows = [];
  // header with ID
  rows.push(
    <tr>
      <td colSpan="4">
        Config ID: {c.id}
      </td>
    </tr>)
  for (var i = 0; i < c.configs.length; i++) {
    var targets = [];
    c.configs[i].targets.forEach((trgt)=> {
      targets.push(
        <Badge pill variant="primary">
          {trgt}
        </Badge>)
    });
    var labels = [];
    console.log(c.configs)
    for (var k in c.configs[i].labels) {
      labels.push(
        <Badge pill variant="secondary">
          {k} = {c.configs[i].labels[k]}
        </Badge>)
    }
    rows.push(
      <tr>
        <td> {i+1} </td>
        <td> {targets} </td>
        <td> {labels} </td>
        <td>
          <Button variant="danger"> Delete </Button>
          <Button variant="warning"> Edit </Button>
        </td>
      </tr>)
  }
  return rows;
}

function buildTable(c) {
  var subTables = [];
  c.forEach((conf) => {
    console.log(conf);
    subTables.push(buildSubTable(conf))
  })
  return <Table striped bordered hover>
    <thead>
      <tr>
        <th>#</th>
        <th>Targets</th>
        <th>Labels</th>
        <th>Actions</th>
      </tr>
    </thead>
    <tbody>
      {subTables}
    </tbody>
  </Table>;
}

function App() {
  var tab = buildTable(configs);

  return (
    <div className="App">
      <Navbar bg="light" expand="lg">
        <Navbar.Brand href="#home">FileSD API</Navbar.Brand>
        <Navbar.Toggle aria-controls="basic-navbar-nav" />
        <Navbar.Collapse id="basic-navbar-nav">
          <Nav className="mr-auto">
            <Nav.Link href="#link">Add</Nav.Link>
            <Nav.Link href="#home">Scrape configs</Nav.Link>
          </Nav>
        </Navbar.Collapse>
      </Navbar>
      <Container>
        <Row>
          {tab}
        </Row>
      </Container>
    </div>
  );
}

export default App;
