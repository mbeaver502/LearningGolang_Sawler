import React, { Component } from 'react';

export default class AppContent extends Component {

    constructor(props) {
        super(props);
        this.handlePostChange = this.handlePostChange.bind(this);
    }

    fetchList = () => {
        fetch('https://jsonplaceholder.typicode.com/posts')
            .then((response) => response.json())
            .then(json => {
                this.handlePostChange(json);
            });
    }

    clickedItem = (x) => {
        console.log("clicked", x);
    }

    handlePostChange = (posts) => {
        this.props.handlePostChange(posts);
    }

    render() {
        return (
            <div>
                This is the content.

                <br />

                <button onClick={this.fetchList} className='btn btn-primary'>Fetch Data</button>

                <hr />

                <p>posts is {this.props.posts.length} items long</p>

                <ul>
                    {this.props.posts.map((c) => (
                        <li key={c.id}>
                            <a href="#!" onClick={() => this.clickedItem(c.id)}>
                                {c.title}
                            </a>
                        </li>
                    ))}
                </ul>
            </div>
        );
    }
}