import React, {useRef, useState} from 'react';
import {Commitable} from "../util/components.jsx";
import {useAppState} from "../util/state.js";
import {Label} from "../ui/Label.jsx";
import {Input} from "../ui/Input.jsx";

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
                <Input
                    type="file"
                    ref={fileInputRef}
                    onChange={handleChange}
                    accept={accept}
                    multiple={multiple}
                />
                <p className="text-sm">{helpText}</p>
            </>
        } />
    );
};
