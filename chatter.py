import asyncio
import nats
import speech_recognition as sr


async def main():
    r = sr.Recognizer()

    with sr.Microphone() as source:
        r.adjust_for_ambient_noise(source)

        print("Please say something")

        audio = r.listen(source)

        try:

            text = r.recognize_google(audio)
            print("Message : " + text)

            nc = await nats.connect("nats://localhost:4222")

            await nc.publish("send_message",  bytes(text,"utf-8"))

            await nc.close()
        
        except Exception as e:
        
            print("Error: " + str(e))

    
if __name__ == "__main__":
   asyncio.run(main())