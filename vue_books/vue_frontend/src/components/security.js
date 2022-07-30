import { store } from './store.js';
import router from './../router/index.js';

let Security = {
    // make sure user is authenticated
    requireToken: function () {
        if (store.token === "") {
            router.push('/login');
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
    },

    // validate user's token
    checkToken: function () {
        if (store.token !== "") {
            const payload = {
                token: store.token
            }

            const headers = new Headers();
            headers.append("Content-Type", "application/json");

            let requestOptions = {
                method: "POST",
                body: JSON.stringify(payload),
                headers: headers,
            }

            fetch(process.env.VUE_APP_API_URL + "/validate-token", requestOptions)
                .then((response) => response.json())
                .then((response) => {
                    if (response.error) {
                        console.log(response.error);
                    } else {
                        if (!response.data) {
                            store.token = "";
                            store.user = {};
                            document.cookie = "_site_data=; Path=/; SameSite=strict; Secure; Expires=Thu, 01 Jan 1970 00:00:01 GMT;";
                        }
                    }
                });
        }
    }
}

export default Security;