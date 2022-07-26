import React, { Component, Fragment } from 'react';
import './EditMovie.css';
import Input from './form-components/Input';
import TextArea from './form-components/TextArea';
import Select from './form-components/Select';

export default class EditMovie extends Component {
    state = {
        movie: {},
        isLoaded: false,
        error: null,
    }

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
        }

        this.handleChange = this.handleChange.bind(this);
        this.handleSubmit = this.handleSubmit.bind(this);
    }

    handleSubmit = (evt) => {
        console.log("form was submitted");
        evt.preventDefault();
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

    componentDidMount() {

    }

    render() {
        let { movie } = this.state
        return (
            <Fragment>
                <h2>Add / Edit Movie</h2>
                <hr />

                <form onSubmit={this.handleSubmit}>
                    <input type='hidden' name='id' id='id' value={movie.id} onChange={this.handleChange} />

                    <Input
                        title={'Title'}
                        type={'text'}
                        name={'title'}
                        value={movie.title}
                        handleChange={this.handleChange}
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

                    <button type='submit' className='btn btn-primary'>Save</button>
                </form>

                <div className='mt-3'>
                    <pre>{JSON.stringify(this.state, null, 3)}</pre>
                </div>
            </Fragment>
        );
    }
}