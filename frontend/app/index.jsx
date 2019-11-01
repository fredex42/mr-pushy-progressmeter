import React from 'react';
import {render} from 'react-dom';
import {BrowserRouter, Link, Route, Switch, Redirect} from 'react-router-dom';

class App extends React.Component {
    render() {
        return <h1>Hello world!</h1>
    }
}

render(<BrowserRouter root="/"><App/></BrowserRouter>, document.getElementById('app'));