package lang

type LanguageText struct {
	Intro string
}

func NewLanguageText() *LanguageText {
	return &LanguageText{
		Intro: `Welcome to SwiftSwap! 🚀

The seamless way to swap your crypto! 💱

Here at SwiftSwipe, we aim to provide you with a fast, secure, and user-friendly experience for all your cryptocurrency swapping needs. Whether you're a seasoned trader or just getting started, our bot is here to assist you every step of the way.

🔒 Secure Transactions: Your security is our top priority. We use cutting-edge encryption to ensure your funds are safe.
⚡ Fast Swaps: Enjoy lightning-fast transactions with minimal fees.
🤝 24/7 Support: Our dedicated support team is always here to help you with any questions or issues.

To get started, simply type /swap and follow the instructions.

Happy swapping with SwiftSwipe! 🚀`,
	}
}
