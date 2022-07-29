import { store } from './store.js';
import router from './../router/index.js';

let Security = {
    // make sure user is authenticated
    requireToken: function () {
        if (store.token === "") {
            router.push('/');
            return false;
        }
    },

    // create request options and send back
    requestOptions: function (payload) {
        const headers = new Headers();
        headers.append("Content-Type", "application/json");
        headers.append("Authorization", "Bearer " + store.token);

        return {
            body: JSON.stringify(payload),
            method: "POST",
            headers: headers,
        };
    }
}

export default Security;