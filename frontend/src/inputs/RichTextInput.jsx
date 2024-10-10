import React, { useState, useEffect } from 'react';
import {Commitable} from "../util/components.jsx";
import {useAppState} from "../util/state.js";

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
                <label className="text-lg font-bold">{label}</label>
                <b>TODO implement rich text!!</b>
                {initialValue}
                <p className="text-sm">{helpText}</p>
            </>
        } />
    );
};
