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
                    className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium file:text-foreground placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                />
                <p className="text-sm">{helpText}</p>
            </>
        } />
    );
};
