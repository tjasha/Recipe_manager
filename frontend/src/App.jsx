// import { useState } from 'react'
// import reactLogo from './assets/react.svg'
// import viteLogo from '/vite.svg'
// import './App.css'
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

  // return (
  //   <>
  //     <div>
  //       <a href="https://vite.dev" target="_blank">
  //         <img src={viteLogo} className="logo" alt="Vite logo" />
  //       </a>
  //       <a href="https://react.dev" target="_blank">
  //         <img src={reactLogo} className="logo react" alt="React logo" />
  //       </a>
  //     </div>
  //     <h1>Vite + React</h1>
  //     <div className="card">
  //       <button onClick={() => setCount((count) => count + 1)}>
  //         count is {count}
  //       </button>
  //       <p>
  //         Edit <code>src/App.jsx</code> and save to test HMR
  //       </p>
  //     </div>
  //     <p className="read-the-docs">
  //       Click on the Vite and React logos to learn more
  //     </p>
  //   </>
  // )

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
