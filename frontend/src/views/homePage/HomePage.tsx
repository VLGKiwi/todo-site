import { ListTask } from "@/module/ListTask/ListTasks"
import { ModalCreateTask } from "@/ui/ModalCreateTask/ModalCreateTask"

export const HomePage = () => {
	return (
		<div>
			<ModalCreateTask isOpen={true} />
			<ListTask />
		</div>
	)
}
