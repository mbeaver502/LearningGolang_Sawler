import React, { useEffect, Fragment } from 'react';
import MoviesFunc from './MoviesFunc';

function AdminFunc(props) {
    useEffect(() => {
        if (props.jwt === "") {
            props.history.push({
                pathname: "/login",
            });

            return;
        }
    }, [props.jwt, props.history]);

    return (
        <Fragment>
            <MoviesFunc
                title="Manage Catalog"
                path="/admin/movie/"
            />
        </Fragment>
    );
}

export default AdminFunc;