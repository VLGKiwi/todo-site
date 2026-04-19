import { create } from "zustand";

type CreateTaskModalState = {
	isOpen: boolean;
	open: () => void;
	close: () => void;
	toggle: () => void;
}

export const useCreateTaskModalStore = create<CreateTaskModalState>((set) => ({
	isOpen: false,
	open: () => set({ isOpen: true }),
	close: () => set({ isOpen: false }),
	toggle: () => set((s) => ({ isOpen: !s.isOpen })),
}))
