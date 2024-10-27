import React, {useRef, useState} from 'react';
import {Commitable} from "../util/components.jsx";
import {useAppState} from "../util/state.js";

export const FileInput = ({ label, helpText, accept, multiple }) => {
    const fileInputRef = useRef(null);
    const [files, setFiles] = useState([]);
    const {sendInput} = useAppState();

    const handleChange = (e) => {
        const files = Array.from(e.target.files);
        setFiles(files);
        console.log('TODO send files', files)
        return true;
    };
    const commit = () => {
        console.log('TODO upload files to the backend', files)
        console.log('TODO notify the script that a file has been uploaded')
        sendInput(files.map(file => file.name))
        return true;
    }

    return (
        <Commitable onCommit={commit} content={
            <>
                <Label>{label}</Label>
                <input
                    type="file"
                    ref={fileInputRef}
                    className="border-gray-900 border-2 outline-2 outline-amber-600 px-4 py-2"
                    onChange={handleChange}
                    accept={accept}
                    multiple={multiple}
                />
                <p className="text-sm">{helpText}</p>
            </>
        } />
    );
};
