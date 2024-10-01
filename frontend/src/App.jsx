import './App.css'
import {create} from 'zustand'
import {useEffect, useState} from "react";
import {marked} from "marked";

const backend = import.meta.env.VITE_BACKEND_URL || 'localhost:8181'

const useAppState = create((set, get) => ({
    pages: [],
    actions: [],
    currentAction: null,
    cards: [],
    history: [],
    startAction: (name) => {
        const msg = {type: 'start', data: name}
        set((state) => ({...state, history: [...state.history, msg], currentAction: name}))
        get().socket.send(JSON.stringify(msg))
    },
    sendInput: (value) => {
        const msg = {type: 'input', data: value}
        set((state) => ({...state, history: [...state.history, msg]}))
        get().socket.send(JSON.stringify(msg))
    }
}))


console.log('Backend URL:', backend)
const displayable = {
    'display': ({text, level}) => {
        const Tag = `h${level}`
        const textStyle = `text-2xl font-bold`
        return <Tag className={textStyle}>{text}</Tag>
    },
    'textInput': ({label, helpText, placeholder, required}) => {
        const {sendInput} = useAppState();
        const [value, setValue] = useState('')
        return (
            <div className={"flex flex-col"}>
                <label className={"text-lg font-bold"}>{label}</label>
                <input className={"border-gray-900 border-2 outline-2 outline-amber-600 px-4 py-2"}
                       onChange={(e) => setValue(e.target.value)} value={value}
                       type={"text"} placeholder={placeholder} required={required}/>
                <p className={"text-sm"}>{helpText}</p>
                <button className={"border-2 border-gray-900 px-4 py-2 rounded hover:bg-amber-600 transition-all"}
                        onClick={() => sendInput(value)}>Submit
                </button>
            </div>
        )
    },
    booleanInput: ({label, helpText, placeholder, required}) => {
        const {sendInput} = useAppState();
        const [value, setValue] = useState(false)
        return (
            <div className={"flex flex-col"}>
                <label className={"text-lg font-bold"}>{label}</label>
                <input type={"radio"} value={value} onClick={() => setValue(true)}/>
                <p className={"text-sm"}>{helpText}</p>
                <button className={"border-2 border-gray-900 px-4 py-2 rounded hover:bg-amber-600 transition-all"}
                        onClick={() => sendInput(value)}>Submit
                </button>
            </div>)
    },
    markdown: ({content}) => (<div className={"prose"} dangerouslySetInnerHTML={{__html: marked(content)}}/>),
}

function setupWebSocket() {
    const socket = new WebSocket(`ws://${backend}`);

    socket.onopen = () => {
        console.log('WebSocket connection established');
        socket.send('{"type": "ping"}');
    };

    socket.onmessage = (event) => {
        try {
            const d = JSON.parse(event.data);
            const {type, data} = d;
            // pages/actions just yeet their state into the store directly
            if (type === 'pages' || type === 'actions') {
                useAppState.setState((state) => ({...state, [type]: data}))
            }
            if (type in displayable) {
                useAppState.setState((state) => ({...state, cards: [...state.cards, {type, data}]}))
            }
            useAppState.setState((state) => ({...state, history: [...state.history, d]}))


        } catch (e) {
            console.error('unparsable message', event.data)
        }
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
            <div className={"flex flex-col gap-2 pb-6"}>
                {app.cards.map((card, i) => {
                    const Displayable = displayable[card.type]
                    return <Displayable key={i} {...card.data} />
                })}
            </div>

            {app.currentAction ? null : (<>
                <h1 className={"text-4xl font-bold"}>Actions</h1>
                <div className={"flex gap-2 pb-6"}>
                    {app.actions && app.actions.map((action, i) => (
                        <button
                            className={"border-2 border-gray-900 px-4 py-2 rounded hover:bg-amber-600 transition-all"}
                            key={action.name} onClick={() => app.startAction(action.name)}>{action.name}</button>
                    ))}
                </div>
            </>)}

            <details className="group border border-gray-200 rounded-lg shadow-sm">
                <summary
                    className="flex justify-between items-center w-full px-4 py-2 text-left text-gray-700 font-medium cursor-pointer focus:outline-none">
                    App State
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
                    <pre>{JSON.stringify(app, null, "  ")}</pre>
                </div>
            </details>

        </div>
    )
}

export default App
