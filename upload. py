import os
import re
import time
import threading
import queue
import subprocess
import requests
import telebot
from telebot import types
from playwright.sync_api import sync_playwright

# ═══════════════════════════════════════
# 🔑 Apni Details
# ═══════════════════════════════════════
TOKEN = "8217681915:AAFhhstt5r3cY1QHaUwdSxn5tlyAaKki_pQ"
CHAT_ID = 6288093086
bot = telebot.TeleBot(TOKEN)

# 🔄 Queue + State
task_queue = queue.Queue()
is_working = False
worker_lock = threading.Lock()
user_context = {
    "state": "IDLE",
    "number": None,
    "otp": None,
    "pending_link": None,
}

BROWSER_ARGS = [
    "--disable-gpu", "--no-sandbox",
    "--disable-dev-shm-usage", "--single-process"
]

YOUTUBE_DOMAINS = ["youtube.com", "youtu.be", "youtube-nocookie.com"]

def is_youtube(link):
    return any(d in link for d in YOUTUBE_DOMAINS)

def safe_filename(title):
    title = re.sub(r'[\\/*?:"<>|]', '', title)
    title = title.strip().replace(' ', '_')
    return title[:80]

def msg(text, **kwargs):
    bot.send_message(CHAT_ID, text, parse_mode="Markdown", **kwargs)

def divider():
    return "─────────────────────"

# ═══════════════════════════════════════
# 📸 Screenshot
# ═══════════════════════════════════════
def take_screenshot(page, caption="📸"):
    try:
        page.screenshot(path="s.png")
        with open("s.png", 'rb') as f:
            bot.send_photo(CHAT_ID, f, caption=caption)
        os.remove("s.png")
    except:
        pass

# ═══════════════════════════════════════
# 🔑 Login
# ═══════════════════════════════════════
def do_login(page, context):
    msg(
        f"╔══════════════════════╗\n"
        f"║   🔐  *LOGIN REQUIRED*   ║\n"
        f"╚══════════════════════╝\n\n"
        f"📱 Apna Jazz number bhejein\n"
        f"Format: `03XXXXXXXXX`"
    )
    user_context["state"] = "WAITING_FOR_NUMBER"

    for _ in range(300):
        if user_context["state"] == "NUMBER_RECEIVED": break
        time.sleep(1)
    else:
        msg("⏰ *Timeout!* Number nahi aaya. Task cancel.")
        return False

    page.locator("#msisdn").fill(user_context["number"])
    time.sleep(1)
    page.locator("#signinbtn").first.click()
    time.sleep(3)
    take_screenshot(page, "📱 Number submit kiya")

    msg(
        f"✅ Number accept hua!\n\n"
        f"🔢 *OTP bhejein* jo aapke\n"
        f"Jazz number pe aaya hai:"
    )
    user_context["state"] = "WAITING_FOR_OTP"

    for _ in range(300):
        if user_context["state"] == "OTP_RECEIVED": break
        time.sleep(1)
    else:
        msg("⏰ *Timeout!* OTP nahi aaya. Task cancel.")
        return False

    for i, digit in enumerate(user_context["otp"].strip()[:6], 1):
        try:
            f = page.locator(f"//input[@aria-label='Digit {i}']")
            if f.is_visible():
                f.fill(digit)
                time.sleep(0.2)
        except: pass

    time.sleep(5)
    take_screenshot(page, "🔢 OTP submit kiya")
    context.storage_state(path="state.json")
    msg(
        f"┌─────────────────────┐\n"
        f"│  ✅  *LOGIN SUCCESSFUL* │\n"
        f"└─────────────────────┘\n\n"
        f"🍪 Session save ho gayi!\n"
        f"Aab link bhejein. 🚀"
    )
    user_context["state"] = "IDLE"
    return True

# ═══════════════════════════════════════
# 🔍 Login Check
# ═══════════════════════════════════════
def check_login_status():
    msg(
        f"🔍 *Jazz Drive* login check\n"
        f"ho raha hai, please wait..."
    )
    with sync_playwright() as p:
        browser = p.chromium.launch(headless=True, args=BROWSER_ARGS)
        ctx = browser.new_context(
            viewport={'width': 1280, 'height': 720},
            storage_state="state.json" if os.path.exists("state.json") else None
        )
        page = ctx.new_page()
        try:
            page.goto("https://cloud.jazzdrive.com.pk/", wait_until="networkidle", timeout=90000)
            time.sleep(3)
            if page.locator("#msisdn").is_visible():
                msg("⚠️ *Session expire ho gayi!*\nLogin karte hain...")
                do_login(page, ctx)
            else:
                msg(
                    f"╔══════════════════════╗\n"
                    f"║  ✅  *LOGIN VALID HAI!*  ║\n"
                    f"╚══════════════════════╝\n\n"
                    f"🚀 Link bhejein — ready hoon!"
                )
        except Exception as e:
            msg(f"❌ Login check error:\n`{str(e)[:150]}`")
        finally:
            browser.close()

# ═══════════════════════════════════════
# 📺 YouTube Quality Keyboard
# ═══════════════════════════════════════
def ask_quality(link):
    user_context["pending_link"] = link
    user_context["state"] = "WAITING_FOR_QUALITY"

    markup = types.ReplyKeyboardMarkup(resize_keyboard=True, one_time_keyboard=True)
    markup.row("🎯 360p", "📱 480p")
    markup.row("💻 720p", "🖥️ 1080p")
    markup.row("⭐ Best Quality")

    msg(
        f"╔══════════════════════╗\n"
        f"║  📺  *YOUTUBE DETECTED*  ║\n"
        f"╚══════════
