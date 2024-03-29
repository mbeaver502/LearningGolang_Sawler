import React, { Component, Fragment } from 'react';
import './EditMovie.css';
import Input from './form-components/Input';
import TextArea from './form-components/TextArea';
import Select from './form-components/Select';
import Alert from './ui-components/Alert';
import { Link } from 'react-router-dom';
import { confirmAlert } from 'react-confirm-alert';
import 'react-confirm-alert/src/react-confirm-alert.css';

export default class EditMovie extends Component {
    constructor(props) {
        super(props);
        this.state = {
            movie: {
                id: 0,
                title: "",
                release_date: "",
                runtime: "",
                mpaa_rating: "",
                rating: "",
                description: "",

            },
            mpaaOptions: [
                { id: 'G', value: 'G' },
                { id: 'PG', value: 'PG' },
                { id: 'PG13', value: 'PG13' },
                { id: 'R', value: 'R' },
                { id: 'NC17', value: 'NC17' },
            ],
            isLoaded: false,
            error: null,
            errors: [],
            alert: {
                type: "d-none",
                message: "",
            },
        }

        this.handleChange = this.handleChange.bind(this);
        this.handleSubmit = this.handleSubmit.bind(this);
    }

    handleSubmit = (evt) => {
        evt.preventDefault();

        // client-side validation
        let errors = [];
        if (this.state.movie.title === '') {
            errors.push('title');
        }

        this.setState({ errors: errors });

        if (errors.length > 0) {
            return false;
        }

        const data = new FormData(evt.target);
        const payload = Object.fromEntries(data.entries());
        const headers = new Headers();

        headers.append("Content-Type", "application/json");
        headers.append("Authorization", "Bearer " + this.props.jwt);

        const requestOptions = {
            method: 'POST',
            body: JSON.stringify(payload),
            headers: headers,
        }

        fetch("http://localhost:4000/v1/admin/editmovie", requestOptions)
            .then(response => response.json())
            .then(data => {
                if (data.error) {
                    this.setState({
                        alert: {
                            type: "alert-danger",
                            message: data.error.message,
                        }
                    });
                } else {
                    this.props.history.push({
                        pathname: "/admin",
                    });
                }
            });
    }

    handleChange = (evt) => {
        let value = evt.target.value;
        let name = evt.target.name;

        this.setState((prevState) => ({
            movie: {
                ...prevState.movie,
                [name]: value,
            }
        }));
    }

    hasError(key) {
        return this.state.errors.indexOf(key) !== -1;
    }

    componentDidMount() {
        if (this.props.jwt === "") {
            this.props.history.push({
                pathname: "/login",
            });

            return;
        }

        const id = this.props.match.params.id;

        // if id > 0, that means the movie already exists
        if (id > 0) {
            fetch("http://localhost:4000/v1/movie/" + id)
                .then((response) => {
                    if (response.status !== 200) {
                        let err = Error;
                        err.message = "Invalid response code";
                        this.setState({
                            isLoaded: true,
                            error: err,
                        });
                    }

                    return response.json();
                })
                .then((json) => {
                    const releaseDate = new Date(json.movie.release_date);

                    this.setState(
                        {
                            movie: {
                                id: id,
                                title: json.movie.title,
                                release_date: releaseDate.toISOString().split("T")[0],
                                runtime: json.movie.runtime,
                                mpaa_rating: json.movie.mpaa_rating,
                                description: json.movie.description,
                                rating: json.movie.rating,
                            },
                            isLoaded: true,
                        },
                        (error) => {
                            this.setState({
                                isLoaded: true,
                                error: error,
                            });
                        }
                    );
                });
        } else {
            this.setState({ isLoaded: true });
        }
    }

    confirmDelete = (e) => {
        if (this.props.jwt === "") {
            this.props.history.push({
                pathname: "/login",
            });

            return;
        }

        confirmAlert({
            title: 'Delete movie?',
            message: 'Are you sure?',
            buttons: [
                {
                    label: 'Yes',
                    onClick: () => {
                        const headers = new Headers();
                        headers.append("Content-Type", "application/json");
                        headers.append("Authorization", "Bearer " + this.props.jwt);

                        const requestOptions = {
                            method: 'GET',
                            headers: headers,
                        }

                        fetch("http://localhost:4000/v1/admin/deletemovie/" + this.state.movie.id, requestOptions)
                            .then(response => response.json())
                            .then(data => {
                                if (data.error) {
                                    this.setState({
                                        alert: {
                                            type: "alert-danger",
                                            message: data.error.message,
                                        }
                                    });
                                } else {
                                    this.props.history.push({
                                        pathname: "/admin",
                                    });
                                }
                            });
                    }
                },
                {
                    label: 'No',
                    onClick: () => { }
                }
            ]
        });
    }

    render() {
        let { movie, isLoaded, error, alert } = this.state

        if (error) {
            return (
                <div>
                    Error: {error.message}
                </div>
            );
        } else if (!isLoaded) {
            return (
                <p>Loading...</p>
            );
        } else {
            return (
                <Fragment>
                    <h2>Add / Edit Movie</h2>

                    <Alert
                        alertType={alert.type}
                        alertMessage={alert.message}
                    />

                    <hr />

                    <form onSubmit={this.handleSubmit}>
                        <input type='hidden' name='id' id='id' value={movie.id} onChange={this.handleChange} />

                        <Input
                            title={'Title'}
                            type={'text'}
                            name={'title'}
                            value={movie.title}
                            handleChange={this.handleChange}
                            className={this.hasError('title') ? 'is-invalid' : ''}
                            errorDiv={this.hasError('title') ? 'text-danger' : 'd-none'}
                            errorMsg={'Please enter a title.'}
                        />

                        <Input
                            title={'Release Date'}
                            type={'date'}
                            name={'release_date'}
                            value={movie.release_date}
                            handleChange={this.handleChange}
                        />

                        <Input
                            title={'Run Time'}
                            type={'text'}
                            name={'runtime'}
                            value={movie.runtime}
                            handleChange={this.handleChange}
                        />

                        <Select
                            name={'mpaa_rating'}
                            title={'MPAA Rating'}
                            value={movie.mpaa_rating}
                            handleChange={this.handleChange}
                            options={this.state.mpaaOptions}
                            placeholder={'Choose...'}
                        />

                        <Input
                            title={'Rating'}
                            type={'text'}
                            name={'rating'}
                            value={movie.rating}
                            handleChange={this.handleChange}
                        />

                        <TextArea
                            name={'description'}
                            title={'Description'}
                            rows={'3'}
                            value={movie.description}
                            handleChange={this.handleChange}
                        />

                        <hr />

                        <button className='btn btn-primary'>Save</button>
                        <Link to="/admin" className="btn btn-warning ms-1">Cancel</Link>
                        {movie.id > 0 && (
                            <a
                                href="#!"
                                onClick={() => this.confirmDelete()}
                                className="btn btn-danger ms-1">
                                Delete
                            </a>
                        )}
                    </form>
                </Fragment>
            );
        }
    }
}