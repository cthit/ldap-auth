import React, { Component } from "react";
import { DigitLoading } from "@cthit/react-digit-components";
import LoginPrompt from "./login-prompt";
import { getBackendUrl } from "../../common/environment";
import { WrapCenter } from "./login.style";
import Axios from "axios";

class Login extends Component {
    constructor(props) {
        super(props);

        let params = new URLSearchParams(props.location.search);
        this.state = {
            client_name: "",
            client_description: "",
            client_id: params.get("client_id"),
            loading: true
        };

        Axios.get(
            `${getBackendUrl()}/application?client_id=${this.state.client_id}`
        )
            .then(res =>
                this.setState({
                    client_name: res.data.name,
                    client_description: res.data.description,
                    loading: false
                })
            )
            .catch(error => {
                console.log("Resource not found");
                console.log(error);
            });
    }

    render = () => (
        <WrapCenter>
            {this.state.loading ? (
                <DigitLoading />
            ) : (
                <LoginPrompt
                    clientName={this.state.client_name}
                    description={this.state.client_description}
                    clientId={this.state.client_id}
                />
            )}
        </WrapCenter>
    );
}
export default Login;
