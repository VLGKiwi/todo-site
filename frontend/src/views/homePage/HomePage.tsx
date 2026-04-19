'use client'

import { ListTasks } from "@/module/ListTasks/ListTasks"
import { useCreateTaskModalStore } from "@/stores/createTaskModalStore"
import { Button } from "@/ui/Button/Button"
import { ModalCreateTask } from "@/ui/ModalCreateTask/ModalCreateTask"

export const HomePage = () => {

	const isOpen = useCreateTaskModalStore((s) => s.isOpen)
	const toggle = useCreateTaskModalStore((s) => s.toggle)

	return (
		<div>
			<Button onClick={toggle}>
				Создать задачу
			</Button>
			<ModalCreateTask isOpen={isOpen} />
			<ListTasks />
		</div>
	)
}
