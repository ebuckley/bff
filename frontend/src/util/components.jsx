import React, {useState} from "react";
import {Card, CardContent, CardFooter, CardHeader} from "../ui/Card.jsx";
import {Button} from "../ui/Button.jsx";

export const Commitable = ({onCommit, content}) => {
    const [hasCommitted, setHasCommitted] = useState(false);

    return (
        <Card>
            <CardHeader></CardHeader>
            <CardContent className={"flex flex-col gap-4"}>
                {React.isValidElement(content) ? content : null}
            </CardContent>
            <CardFooter>
                {hasCommitted ? (
                    <p className={"text-sm text-gray-500"}>Submitted</p>
                ) : (
                    <Button
                        onClick={() => {
                            if (onCommit()) {
                                setHasCommitted(true);
                            }
                        }}
                    >
                        Submit
                    </Button>
                )}
            </CardFooter>
        </Card>
    );
};
