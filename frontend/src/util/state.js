import {create} from "zustand";

export const backend = `${window.location.host}${window.location.pathname}/ws`

export const actionName = window.location.pathname.split('/').pop()

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
