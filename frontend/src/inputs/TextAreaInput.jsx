import React, {useState} from 'react';
import {Commitable} from "../util/components.jsx";
import {useAppState} from "../util/state.js";
import {Label} from "../ui/Label.jsx";

export const TextAreaInput = ({label, helpText, placeholder, required}) => {
    const {sendInput} = useAppState();
    const [value, setValue] = useState('');

    const handleChange = (e) => {
        setValue(e.target.value);
    };

    const handleCommit = () => {
        sendInput(value);
        return true;
    };

    return (
        <Commitable onCommit={handleCommit} content={
            <>
                <Label>{label}</Label>
                <textarea
                    className="border-gray-900 border-2 outline-2 outline-amber-600 px-4 py-2"
                    onChange={handleChange}
                    placeholder={placeholder}
                    required={required}
                >
      {value}
    </textarea>
                <p className="text-sm">{helpText}</p>
            </>
        }/>
    );
};
