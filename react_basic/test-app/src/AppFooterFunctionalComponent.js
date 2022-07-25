import React, { Fragment } from 'react';
import './AppFooter.css';

export default function AppFooterFunctionalComponent(props) {
    const currentYear = new Date().getFullYear();
    return (
        <Fragment>
            <hr />
            <p className='footer'>
                Copyright &copy; 2020 - {currentYear} Acme Ltd.
                <br />
                {props.myProperty}
            </p>
        </Fragment>
    );
}