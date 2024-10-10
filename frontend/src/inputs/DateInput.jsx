import React, { useState } from 'react';
import DatePicker from 'react-datepicker'; // You'll need to install this package
import 'react-datepicker/dist/react-datepicker.css';
import {Commitable} from "../util/components.jsx";
import {useAppState} from "../util/state.js";

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
                <label className="text-lg font-bold">{label}</label>
                <DatePicker
                    selected={selectedDate}
                    onChange={handleChange}
                    minDate={new Date(min)}
                    maxDate={new Date(max)}
                    className="border-gray-900 border-2 outline-2 outline-amber-600 px-4 py-2"
                 showMonthYearDropdown/>
                <p className="text-sm">{helpText}</p>
            </>
        } />
    );
};
