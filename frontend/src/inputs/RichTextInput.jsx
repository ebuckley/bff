import React, { useState, useEffect } from 'react';
import {Commitable} from "../util/components.jsx";
import {useAppState} from "../util/state.js";
import {Label} from "../ui/Label.jsx";

export const RichTextInput = ({ label, helpText, initialValue }) => {
    const {sendInput} = useAppState();

    let handleCommit = () => {
        console.log('submitting the rich text state...')
        sendInput(initialValue);
        return true;
    };
    return (
        <Commitable onCommit={handleCommit} content={
            <>
                <Label>{label}</Label>
                <b>TODO implement rich text!!</b>
                {initialValue}
                <p className="text-sm">{helpText}</p>
            </>
        } />
    );
};
