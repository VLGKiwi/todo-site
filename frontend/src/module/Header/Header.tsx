import Image from "next/image"
import Link from "next/link"
import styles from './Header.module.scss'

export const Header = () => {
	return (
		<header className={styles.header}>
			<Link href={'/'}>
				<Image
					src={'./icons/logo.svg'}
					width={200}
					height={64}
					alt="Logo"
				/>
			</Link>
		</header>
	)
}
