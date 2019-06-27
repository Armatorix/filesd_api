import React from 'react';
import { Container, Row, Badge, Table, Navbar, Nav } from 'react-bootstrap'
import './App.css';

var configs = [
  {
    targets: ["123.321.213.321:3213", "123.321.213.321:3213", "123.321.213.321:3213", "123.321.213.321:3213"],
    labels: { "1": "1lab", "2": "2lab" },
  },
  {
    targets: ["123.321.213.321:3213", "123.321.213.321:3213", "123.321.213.321:3213", "123.321.213.321:3213"],
    labels: { "1": "1lab", "2": "2lab" },
  }
]

function App() {
  var rows = [];
  for (var i = 0; i < configs.length; i++) {
    var targets = [];
    for (var j = 0; j < configs[i].targets.length; j++) {
      targets.push(
        <Badge pill variant="primary">
          {configs[i].targets[j]}
        </Badge>)
    }
    var labels = [];
    for (var k in configs[i].labels) {
      labels.push(
        <Badge pill variant="secondary">
          {k} = {configs[i].labels[k]}
        </Badge>)
    }
    rows.push(
      <tr>
        <td> {i} </td>
        <td> {targets} </td>
        <td> {labels} </td>
      </tr>)
  }
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
          <Table striped bordered hover>
            <thead>
              <tr>
                <th>#</th>
                <th>Targets</th>
                <th>Labels</th>
              </tr>
            </thead>
            <tbody>
              {rows}
            </tbody>
          </Table>
        </Row>
      </Container>
    </div>
  );
}

export default App;
