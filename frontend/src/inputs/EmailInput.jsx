import React, { useState } from 'react';
import {Commitable} from "../util/components.jsx";
import {useAppState} from "../util/state.js";
import {Input} from "../ui/Input.jsx";

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

                <Input
                    type="email"
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
