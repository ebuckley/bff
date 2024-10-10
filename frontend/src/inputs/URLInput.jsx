import React, { useState } from 'react';
import {Commitable} from "../util/components.jsx";
import {useAppState} from "../util/state.js";

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
                <label className="text-lg font-bold">{label}</label>
                <input
                    type="url"
                    className="border-gray-900 border-2 outline-2 outline-amber-600 px-4 py-2"
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
