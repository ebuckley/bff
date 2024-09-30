import './App.css'
import {create} from 'zustand'

const backend = import.meta.env.VITE_BACKEND_URL || 'localhost:8181'

const useAppState = create((set) => ({
    pages: null,
    actions: null,
    onMessage: (msg) => {
        console.log('onMessage', msg)
        const {Type, Data} = msg
        set(state => ({
            ...state,
            [Type]: Data
        }))
    },
    startAction: (name) => {
        socket.send(JSON.stringify({type: 'start', name}))
    },
}))


console.log('Backend URL:', backend)


const socket = new WebSocket(`ws://${backend}`);

socket.onopen = () => {
    console.log('WebSocket connection established');
    socket.send('{"type": "ping"}');
};

socket.onmessage = (event) => {
    useAppState.getState().onMessage(JSON.parse(event.data));
    // useAppState.getState().onMessage(JSON.parse(event.data));
};

socket.onclose = () => {
    console.log('WebSocket connection closed');
};

socket.onerror = (error) => {
    console.error('WebSocket error:', error);
};


function App() {

    const app = useAppState(s => s)
    const actions = useAppState(s => s.actions)

    return (
        <div className={"py-6 mx-auto max-w-2xl"}>
            <pre>{JSON.stringify(app, null, "  ")}</pre>

            {actions && actions.map((action, i) => (
                <button className={"border-2 border-gray-900 px-4 py-2 rounded hover:bg-amber-600 transition-all"} key={action.name} onClick={()=> app.startAction(action.name)}>{action.name}</button>
            ))}
        </div>
    )
}

export default App
