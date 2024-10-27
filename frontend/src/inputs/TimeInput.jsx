import React, { useState } from 'react';
import TimePicker from 'react-time-picker'; // You'll need to install this package
import 'react-time-picker/dist/TimePicker.css';
import {Commitable} from "../util/components.jsx";
import {useAppState} from "../util/state.js";
import {Label} from "../ui/Label.jsx";


export const TimeInput = ({ label, helpText, min, max, onCommit }) => {
    const [time, setTime] = useState('12:00');
    const {sendInput} = useAppState();
    const handleChange = (newTime) => {
        setTime(newTime);
    };

    const handleCommit = () => {
        sendInput(time);
        return true
    };

    return (
        <Commitable onCommit={handleCommit} content={
            <>
                <Label>{label}</Label>
                <TimePicker
                    onChange={handleChange}
                    value={time}
                    minTime={min}
                    maxTime={max}
                    disableClock={true} // clock css does not work with the tailwind preset
                    className="border-gray-900 border-2 outline-2 outline-amber-600 px-4 py-2"
                />
                <p className="text-sm">{helpText}</p>
            </>
        } />
    );
};
