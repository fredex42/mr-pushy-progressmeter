import React from 'react';
import {render} from 'react-dom';
import {BrowserRouter, Link, Route, Switch, Redirect} from 'react-router-dom';
import RunsList from "./runslist.jsx";

class App extends React.Component {
    render() {
        return <div>
            <h1>Mr Pushy Monitor</h1>
            <div className="sidebar">
                <h2>Runs</h2>
                <RunsList/>
            </div>
        </div>
    }
}

render(<BrowserRouter root="/"><App/></BrowserRouter>, document.getElementById('app'));