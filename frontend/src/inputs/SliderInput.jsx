import React, { useState } from 'react';
import Slider from 'react-slider';
import {Commitable} from "../util/components.jsx";
import {useAppState} from "../util/state.js"; // You'll need to install this package

export const SliderInput = ({ label, helpText, min, max, step }) => {
    const [value, setValue] = useState(min);
    const {sendInput} = useAppState();

    const handleChange = (newValue) => {
        setValue(newValue);
    };

    const handleCommit = () => {
        sendInput(value);
        return true;
    };

    return (
        <Commitable onCommit={handleCommit} content={
            <>
                <label className="text-lg font-bold">{label}</label>
                <Slider
                    min={min}
                    max={max}
                    step={step}
                    value={value}
                    onChange={handleChange}
                    className="border-1 bg-gray-100 w-full h-8 rounded-full"
                    markClassName="bg-amber-600 w-8 h-8 rounded-full"
                    thumbClassName="border-1 bg-gray-900 w-8 h-8 rounded-full font-bold text-white text-sm flex justify-center items-center"
                    renderThumb={(props, state) => <div {...props}>{state.valueNow}</div>}
                />
                <p>{value}</p>
                <p className="text-sm">{helpText}</p>
            </>
        } />
    );
};
