import React, { Component, Fragment } from 'react';
import Movies from './Movies';

export default class Admin extends Component {

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