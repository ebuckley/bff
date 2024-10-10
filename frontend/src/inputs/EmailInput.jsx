import React, { useState } from 'react';
import {Commitable} from "../util/components.jsx";
import {useAppState} from "../util/state.js";

export const EmailInput = ({ label, helpText, placeholder, required }) => {
    const {sendInput} = useAppState();
    const [value, setValue] = useState('');

    const handleChange = (e) => {
        setValue(e.target.value);
    };

    const handleCommit = () => {
        if (value.match(/^[^\s@]+@[^\s@]+\.[^\s@]+$/)) {
            sendInput(value);
            return true;
        } else {
            alert('Please enter a valid email address');
            return false;
        }
    };

    return (
        <Commitable onCommit={handleCommit} content={
            <>
                <label className="text-lg font-bold">{label}</label>
                <input
                    type="email"
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
