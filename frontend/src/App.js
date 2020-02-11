import React from "react";
import { BrowserRouter, Switch } from "react-router-dom";
import { DigitProviders } from "@cthit/react-digit-components";
import { Route } from "react-router";
import Login from "./use-cases/login";
import NotFound from "./use-cases/not-found";
import "./App.css";

const App = () => {
    return (
        <>
            <DigitProviders>
                <BrowserRouter>
                    <Switch>
                        <Route exact path="/authenticate" component={Login} />
                        } />
                        <Route path="/" component={NotFound} />
                    </Switch>
                </BrowserRouter>
            </DigitProviders>
        </>
    );
};

export default App;
