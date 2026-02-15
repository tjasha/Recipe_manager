import React, { useEffect, useState } from 'react';
import { Container, Navbar, Nav, Table } from 'react-bootstrap';
import 'bootstrap/dist/css/bootstrap.min.css';

function App() {
  const [message, setMessage] = useState("");


    useEffect(() => {
        fetch("/api/hello")
            .then(res => res.json())
            .then(data => setMessage(data.message));
    }, []);

    return (
        <Container>
            <Navbar bg="dark" variant="dark">
                <Navbar.Brand href="#">Recipes</Navbar.Brand>
                <Nav className="me-auto">
                    <Nav.Link href="#">Home</Nav.Link>
                    <Nav.Link href="#">Favourites</Nav.Link>
                    <Nav.Link href="#">My Recipes</Nav.Link>
                    <Nav.Link href="#">Log in</Nav.Link>
                </Nav>
            </Navbar>

            <h1 className="mt-5">Welcome to Recipe manager!</h1>
            <p>{message}</p>

            <h2>Newest Recipes</h2>
            <Table striped bordered hover>
                <thead>
                <tr>
                    <th>Name</th>
                    <th>Ingredients</th>
                </tr>
                </thead>
                <tbody>
                <tr>
                    <td>Pasta bake</td>
                    <td>Pasta, Tomato, Zucchini, Ricotta</td>
                </tr>
                <tr>
                    <td>Mediterranean Lentil Salad</td>
                    <td>Lentils, Avocado, Lettuce, Tomato, Cucumber</td>
                </tr>
                </tbody>
            </Table>
        </Container>
    );

}

export default App;
