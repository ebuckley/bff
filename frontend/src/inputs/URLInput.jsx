import React, { useState } from 'react';
import {Commitable} from "../util/components.jsx";
import {useAppState} from "../util/state.js";
import {Input} from "postcss";

export const URLInput = ({ label, helpText, placeholder, required, onCommit }) => {
    const [value, setValue] = useState('');
    const {sendInput} = useAppState();
    const handleChange = (e) => {
        setValue(e.target.value);
    };

    const handleCommit = () => {
        if (value.match(/^(https?:\/\/)?([\da-z\.-]+)\.([a-z\.]{2,6})([\/\w \.-]*)*\/?$/)) {
            sendInput(value);
            return true;
        } else {
            alert('Please enter a valid URL');
            return false;
        }
    };

    return (
        <Commitable onCommit={handleCommit} content={
            <>
                <Label>{label}</Label>
                <Input
                    type="url"
                    onChange={handleChange}
                    value={value}
                    placeholder={placeholder}
                    required={required}
                />
                <p className="text-sm">{helpText}</p>
            </>
        } />
    );
};
