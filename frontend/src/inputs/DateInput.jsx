import React, { useState } from 'react';
import DatePicker from 'react-datepicker'; // You'll need to install this package
import 'react-datepicker/dist/react-datepicker.css';
import {Commitable} from "../util/components.jsx";
import {useAppState} from "../util/state.js";
import {Label} from "../ui/Label.jsx";

export const DateInput = ({ label, helpText, min, max }) => {
    const [selectedDate, setSelectedDate] = useState(null);
    const {sendInput} = useAppState();
    const handleChange = (date) => {
        setSelectedDate(date);
    };

    const handleCommit = () => {
        if (selectedDate) {
            sendInput(selectedDate.toISOString().split('T')[0])
            return true;
        }
        return false;
    };

    return (
        <Commitable onCommit={handleCommit} content={
            <>
                <Label>{label}</Label>
                <DatePicker
                    selected={selectedDate}
                    onChange={handleChange}
                    minDate={new Date(min)}
                    maxDate={new Date(max)}
                    className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium file:text-foreground placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                 showMonthYearDropdown/>
                <p className="text-sm">{helpText}</p>
            </>
        } />
    );
};
