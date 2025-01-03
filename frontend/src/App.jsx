import './App.css'
import React, {useEffect, useState} from "react";
import {marked} from "marked";
import SyntaxHighlighter from 'react-syntax-highlighter';
import {atomDark} from "react-syntax-highlighter/src/styles/prism/index.js";
import {Commitable} from "./util/components.jsx";
import {FileInput} from "./inputs/FileInput.jsx";
import {EmailInput} from "./inputs/EmailInput.jsx";
import {DateInput} from "./inputs/DateInput.jsx";
import {RichTextInput} from "./inputs/RichTextInput.jsx";
import {URLInput} from "./inputs/URLInput.jsx";
import {TimeInput} from "./inputs/TimeInput.jsx";
import {SliderInput} from "./inputs/SliderInput.jsx";
import {backend, useAppState, actionName} from "./util/state.js";
import {TextAreaInput} from "./inputs/TextAreaInput.jsx";
import {Input} from "./ui/Input.jsx";
import {Switch} from "./ui/Switch.jsx";
import {Label} from "./ui/Label.jsx";
import {Card, CardContent, CardHeader} from "./ui/Card.jsx";


console.log('Backend URL:', backend)
const displayable = {
    'image': ({url, alt}) => <img src={url} alt={alt}/>,
    'display': ({text, level}) => {
        const Tag = `h${level}`
        const textStyle = `text-2xl font-bold`
        return <Tag className={textStyle}>{text}</Tag>
    },
    'numberInput': ({label, helpText, placeholder, required}) => {
        const {sendInput} = useAppState();

        const [value, setValue] = useState('')

        const commitSend = () => {
            sendInput(value)
            return true
        }
        return (
            <Commitable onCommit={commitSend} content={<>
                <Label>{label}</Label>
                <Input
                       onChange={(e) => setValue(e.target.value)} value={value}
                       type={"number"} placeholder={placeholder} required={required}/>
                <p className={"text-sm"}>{helpText}</p>
            </>}/>
        )
    },
    'textInput': ({label, helpText, placeholder, required}) => {
        const {sendInput} = useAppState();

        const [value, setValue] = useState('')

        const commitSend = () => {
            sendInput(value)
            return true
        }
        return (
            <Commitable onCommit={commitSend} content={<>
                <Label >{label}</Label>
                <Input
                       onChange={(e) => setValue(e.target.value)} value={value}
                       type={"text"} placeholder={placeholder} required={required}/>
                <p className={"text-sm"}>{helpText}</p>
            </>}/>
        )
    },
    booleanInput: ({label, helpText, placeholder, required}) => {
        const {sendInput} = useAppState();
        const [value, setValue] = useState(false)
        return (
            <Commitable onCommit={() => {
                sendInput(value)
                return true;
            }} content={<div className={"flex flex-col gap-2"}>
                <Label>{label}</Label>
                <Switch value={value} onClick={() => setValue(true)}/>
                <p className={"text-sm"}>{helpText}</p>
            </div>}/>)
    },
    markdown: ({content}) => (<div className={"prose"} dangerouslySetInnerHTML={{__html: marked(content)}}/>),

    'link': ({ text, url, type }) => {
        const baseStyle = "px-4 py-2 rounded-md text-white";
        const typeStyles = {
            default: "bg-blue-500 hover:bg-blue-600",
            primary: "bg-green-500 hover:bg-green-600",
            danger: "bg-red-500 hover:bg-red-600",
        };
        const buttonStyle = `${baseStyle} ${typeStyles[type] || typeStyles.default}`;

        return (
            <a href={url} className={buttonStyle} target="_blank" rel="noopener noreferrer">
                {text}
            </a>
        );
    },

    'html': ({ content }) => (
        <div dangerouslySetInnerHTML={{ __html: content }} />
    ),

    'code': ({ code, language }) => (
        <SyntaxHighlighter language={language || 'text'} style={atomDark}>
            {code}
        </SyntaxHighlighter>
    ),

    'metadata': ({ items, layout }) => {
        const renderItems = () => {
            return items.map((item, index) => (
                <div key={index} className="mb-2">
                    <span className="font-bold">{item.label}: </span>
                    <span>{item.value}</span>
                </div>
            ));
        };

        const layoutStyles = {
            default: "bg-white p-4",
            table: "table-auto",
        };

        if (layout === 'table') {
            return (
                <table className={layoutStyles.table}>
                    <tbody>
                    {items.map((item, index) => (
                        <tr key={index}>
                            <td className="font-bold pr-4">{item.label}</td>
                            <td>{item.value}</td>
                        </tr>
                    ))}
                    </tbody>
                </table>
            );
        }

        return (
            <Card className={layoutStyles[layout] || layoutStyles.default}>
                <CardHeader></CardHeader>
                <CardContent>
                    {renderItems()}
                </CardContent>
            </Card>
        );
    },
    'emailInput': EmailInput,
    'sliderInput': SliderInput,
    'dateInput': DateInput,
    'richTextInput': RichTextInput,
    'urlInput': URLInput,
    'timeInput': TimeInput,
    'fileInput': FileInput,
    'textAreaInput': TextAreaInput,
}

function setupWebSocket() {
    const socket = new WebSocket(`${window.location.protocol === 'https:' ? 'wss' : 'ws'}://${backend}`);

    socket.onopen = () => {
        console.log('WebSocket connection established');
        socket.send('{"type": "ping"}');
        const s = useAppState.getState()
        s.startAction(actionName)
    };

    socket.onmessage = (event) => {
        let d = {};
        try {
            d = JSON.parse(event.data);
        } catch (e) {
            console.error('unparsable message', event.data)
            return;
        }
        const {type, data} = d;
        // pages/actions just yeet their state into the store directly
        if (type === 'pages' || type === 'actions') {
            useAppState.setState((state) => ({...state, [type]: data}))
        }
        if (type in displayable) {
            useAppState.setState((state) => ({...state, cards: [...state.cards, {type, data}]}))
        }
        if (type === 'done') {
            // todo send something into state for rendering that this is done ta-da
            useAppState.setState((state) => ({...state, currentAction: null}))
        }

        // also append the message to the global history of messages
        useAppState.setState((state) => ({...state, history: [...state.history, d]}))
    };

    socket.onclose = () => {
        console.log('WebSocket connection closed');
        // todo start retry connection?
    };

    socket.onerror = (error) => {
        console.error('WebSocket error:', error);
    };

    useAppState.setState((state) => ({...state, socket}))
    return socket
}

function App() {
    const app = useAppState()
    const errors = app.history.filter((msg) => msg.type === 'error').map((msg,i) => (
    <div key={i} className="py-6 px-3 bg-red-400 color-red-900 rounded border-2 border-red-900">
        <span className={"font-bold pr-1"}> Error </span> {msg.data}
    </div>))
    useEffect(() => {
        const socket = setupWebSocket()
        return () => {
            if (socket) {
                socket.close()
            }
        }
    }, [])
    return (
        <div className={"py-6 mx-auto max-w-2xl"}>
            <div className="flex flex-col gap-2 pb-3">
                {errors}
            </div>
            <div className={"flex flex-col gap-2 pb-6"}>
                {app.cards.map((card, i) => {
                    const Displayable = displayable[card.type]
                    return <Displayable key={i} {...card.data} />
                })}
            </div>

            <details className="group border border-gray-200 rounded-lg shadow-sm">
                <summary
                    className="flex justify-between items-center w-full px-4 py-2 text-left text-gray-700 font-medium cursor-pointer focus:outline-none">
                    Debug App State
                    <svg
                        className="w-5 h-5 text-gray-500 transition-transform duration-200 group-open:rotate-180"
                        xmlns="http://www.w3.org/2000/svg"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                    >
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7"/>
                    </svg>
                </summary>
                <div className="px-4 pb-4 text-sm text-gray-600">
                    <div>Backend: {backend}</div>
                    <pre>{JSON.stringify(app, null, "  ")}</pre>
                </div>
            </details>

        </div>
    )
}

export default App
