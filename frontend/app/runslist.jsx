import React from 'react';

class RunsList extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            loading: false,
            runsList: [],
            lastError: null
        }
    }

    componentDidMount() {
        this.loadData();
    }

    loadData() {
        this.setState({loading:true}, ()=>{
            fetch("/api/job/list").then(result=>{
                result.json().then(data=> {
                    console.log("Got ", data);
                    this.setState({
                        loading: false,
                        runsList: data.jobs
                    })
                })
            }).catch(error=>{
                console.error(error);
                this.setState({loading: false, lastError: error})
            })
        })
    }

    render() {
        if(this.state.loading){
            return <span className="information">Loading...</span>
        } else if(this.state.lastError) {
            return <span className="error">{this.state.lastError}</span>
        } else {
            if(this.state.runsList && this.state.runsList.length>0) {
                return <ul>{
                    this.state.runsList.map(runEntry =>
                        <li className="selector">
                            Started at: {runEntry.timestamp}<br/>
                            Job ID: {runEntry.jobId}<br/>
                            Total size: {runEntry.uploadsCount}<br/>
                            Status: {runEntry.status}<br/>
                        </li>)
                }</ul>
            } else {
                return <span className="information">No runs registered yet</span>
            }
        }
    }
}

export default RunsList;