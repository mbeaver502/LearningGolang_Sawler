import React, { Component, Fragment } from 'react';
import Movies from './Movies';

export default class Admin extends Component {
    componentDidMount() {
        if (this.props.jwt === "") {
            this.props.history.push({
                pathname: "/login",
            });

            return;
        }
    }

    render() {
        return (
            <Fragment>
                <Movies
                    title="Manage Catalog"
                    path="/admin/movie/"
                />
            </Fragment>
        );
    }
}