import React, {useState} from 'react';
import {Commitable} from "../util/components.jsx";
import {useAppState} from "../util/state.js";

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
                <label className="text-lg font-bold">{label}</label>
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
