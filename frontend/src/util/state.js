import {create} from "zustand";

export const backend = import.meta.env.VITE_BACKEND_URL || 'localhost:8181'

export const useAppState = create((set, get) => ({
    pages: [],
    actions: [],
    currentAction: null,
    cards: [],
    history: [],
    startAction: (name) => {
        const msg = {type: 'start', data: name}
        set((state) => ({...state, history: [...state.history, msg], currentAction: name,  cards: []}))
        get().socket.send(JSON.stringify(msg))
    },
    sendInput: (value) => {
        const msg = {type: 'input', data: value}
        set((state) => ({...state, history: [...state.history, msg]}))
        get().socket.send(JSON.stringify(msg))
    }
}))
