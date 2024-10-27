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
                    className="flex min-h-[80px] w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                    onChange={handleChange}
                    placeholder={placeholder}
                    required={required}
                >{value}</textarea>
                <p className="text-sm">{helpText}</p>
            </>
        }/>
    );
};
