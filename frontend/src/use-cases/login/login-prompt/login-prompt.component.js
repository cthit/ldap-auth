import React from "react";
import {
    DigitEditDataCard,
    DigitTextField,
    useDigitToast
} from "@cthit/react-digit-components";
import Axios from "axios";
import { getBackendUrl } from "../../../common/environment";

const LoginPrompt = ({ clientName, description, clientId }) => {
    const [toast] = useDigitToast({
        duration: 3000,
        actionText: "Close",
        actionHandler: () => {}
    });
    return (
        <DigitEditDataCard
            initialValues={{
                cid: "",
                password: ""
            }}
            onSubmit={(values, actions) => {
                actions.setSubmitting(true);

                Axios.post(
                    `${getBackendUrl()}/authenticate?client_id=${clientId}`,
                    {
                        cid: values.cid,
                        password: values.password
                    }
                )
                    .then(res => {
                        window.location.replace(
                            `${res.callback_url}?token=${res.data.token}`
                        );
                    })
                    .catch(error => {
                        console.log(error);
                        toast({ text: "Failed to log in" });
                        actions.setSubmitting(false);
                    });
            }}
            titleText={`Login for ${clientName}`}
            subtitleText={description}
            submitText={"Login"}
            keysOrder={["cid", "password"]}
            keysComponentData={{
                cid: {
                    component: DigitTextField,
                    componentProps: {
                        upperLabel: "CID"
                    }
                },
                password: {
                    component: DigitTextField,
                    componentProps: {
                        upperLabel: "Password",
                        password: true
                    }
                }
            }}
        />
    );
};

export default LoginPrompt;
