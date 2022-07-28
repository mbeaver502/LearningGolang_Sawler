const RegistrationForm = {
    data() {
        return {
            addressSameChecked: true,
        };
    },
    props: ["items"],
    template:
        `<h3>Registration</h3>
        <hr/>
        <form autocomplete="off" class="needs-validation" novalidate>
            <text-input label="First Name" name="first_name" required="required" type="text"></text-input>
            <text-input label="Last Name" name="last_name" required="required" type="text"></text-input>
            <text-input label="Email" name="email" required="required" type="email"></text-input>
            <text-input label="Password" name="password" required="required" type="password"></text-input>
            <select-input label="Favorite Color" name="color" :items="items"></select-input>

            <text-input label="Address" name="address" required="required" type="text"></text-input>
            <text-input label="City Name" name="city" required="required" type="text"></text-input>
            <text-input label="State" name="state" required="required" type="text"></text-input>
            <text-input label="ZIP Code" name="zip_code" required="required" type="text"></text-input>

            <check-input v-on:click="addressSame" label="Mailing Address Same" checked="true" v-model="addressSameChecked"></check-input>

            <div v-if="addressSameChecked === false">
                <div class="mt-3">
                    <text-input label="Mailing Address" name="mailing_address" type="text"></text-input>
                    <text-input label="Mailing City Name" name="mailing_city"  type="text"></text-input>
                    <text-input label="Mailing State" name="mailing_state"  type="text"></text-input>
                    <text-input label="Mailing ZIP Code" name="mailing_zip_code" type="text"></text-input>
                </div>
            </div>

            <check-input label="Agree to ToS" name="agree" required="required" value="1"></check-input>

            <hr/>

            <input class="btn btn-outline-primary" type="submit" value="Register">
        </form>`,
    methods: {
        addressSame() {
            console.log("address same fired");
            if (this.addressSameChecked === true) {
                console.log("was checked on click");
                this.addressSameChecked = false;
            } else {
                console.log("was NOT checked on click");
                this.addressSameChecked = true;
            }
        }
    },
    components: {
        'text-input': TextInput,
        'select-input': SelectInput,
        'check-input': CheckboxInput,
    },
    mounted() {
        // Example starter JavaScript for disabling form submissions if there are invalid fields
        (() => {
            'use strict'

            // Fetch all the forms we want to apply custom Bootstrap validation styles to
            const forms = document.querySelectorAll('.needs-validation')

            // Loop over them and prevent submission
            Array.from(forms).forEach(form => {
                form.addEventListener('submit', event => {
                    if (!form.checkValidity()) {
                        event.preventDefault()
                        event.stopPropagation()
                    }

                    form.classList.add('was-validated')
                }, false)
            })
        })()
    }
}