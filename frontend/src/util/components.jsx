import React, {useState} from "react";

export const Commitable = ({ onCommit, content }) => {
    const [hasCommitted, setHasCommitted] = useState(false);

    return (
        <div className={"flex flex-col"}>
            {React.isValidElement(content) ? content : null}
            {hasCommitted ? (
                <p className={"text-sm text-gray-500"}>Submitted</p>
            ) : (
                <button
                    className={"border-2 border-gray-900 px-4 py-2 rounded hover:bg-amber-600 transition-all"}
                    onClick={() => {
                        if (onCommit()) {
                            setHasCommitted(true);
                        }
                    }}
                >
                    Submit
                </button>
            )}
        </div>
    );
};
