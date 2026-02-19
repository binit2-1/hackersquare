"use client";

import React, { useState, useEffect, useRef } from "react";

const searchBarStyles = `
@keyframes placeholder-slide-up {
  0% {
    transform: translateY(0);
    opacity: 1;
    filter: blur(0px);
  }
  100% {
    transform: translateY(-50%);
    opacity: 0;
    filter: blur(3px);
  }
}

@keyframes placeholder-slide-in {
  0% {
    transform: translateY(50%);
    opacity: 0;
    filter: blur(3px);
  }
  100% {
    transform: translateY(0);
    opacity: 1;
    filter: blur(0px);
  }
}
`;

export function SearchBar({
  placeholders = ["Search...", "Type something...", "What are you looking for?"],
  interval = 3000,
  onChange,
  onSubmit,
  icon,
  className = "",
  inputClassName = "",
  ...props
}) {
  const [currentIndex, setCurrentIndex] = useState(0);
  const [isAnimating, setIsAnimating] = useState(false);
  const [inputValue, setInputValue] = useState("");
  const [isFocused, setIsFocused] = useState(false);
  const inputRef = useRef(null);

  useEffect(() => {
    if (inputValue) return;

    const timer = setInterval(() => {
      setIsAnimating(true);
      setTimeout(() => {
        setCurrentIndex((prev) => (prev + 1) % placeholders.length);
        setIsAnimating(false);
      }, 300);
    }, interval);

    return () => clearInterval(timer);
  }, [placeholders.length, interval, inputValue]);

  // shortcut: Ctrl+K / Cmd+K to focus
  useEffect(() => {
    const handleGlobalKeyDown = (e) => {
      if ((e.ctrlKey || e.metaKey) && e.key === "k") {
        e.preventDefault();
        inputRef.current?.focus();
      }
    };

    window.addEventListener("keydown", handleGlobalKeyDown);
    return () => window.removeEventListener("keydown", handleGlobalKeyDown);
  }, []);

  const handleChange = (e) => {
    setInputValue(e.target.value);
    onChange?.(e);
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    onSubmit?.(inputValue);
  };

  const handleKeyDown = (e) => {
    if (e.key === "Enter") {
      handleSubmit(e);
    }
  };

  const handleClear = () => {
    setInputValue("");
    inputRef.current?.focus();
  };

  return (
    <>
      <style>{searchBarStyles}</style>
      <form onSubmit={handleSubmit} className={`relative w-full max-w-lg ${className}`}>
        <div
          className={`
            relative flex items-center w-full px-1
            bg-white/3 backdrop-blur-2xl
            border border-white/8
            rounded-full
            transition-all duration-300 ease-out
            hover:border-white/12 hover:bg-white/4
            focus-within:border-white/12 focus-within:bg-white/4
            shadow-[0_4px_24px_-8px_rgba(0,0,0,0.3)]
          `}
        >
          <div className="flex items-center justify-center w-12 h-12 text-zinc-500">
            {icon || (
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="18"
                height="18"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round"
              >
                <circle cx="11" cy="11" r="8" />
                <path d="m21 21-4.3-4.3" />
              </svg>
            )}
          </div>

          <div className="relative flex-1 h-12">
            <input
              ref={inputRef}
              type="text"
              value={inputValue}
              onChange={handleChange}
              onFocus={() => setIsFocused(true)}
              onBlur={() => setIsFocused(false)}
              onKeyDown={handleKeyDown}
              className={`
                w-full h-full bg-transparent
                text-zinc-200 text-[15px]
                outline-none
                placeholder-transparent
                pr-4
                ${inputClassName}
              `}
              {...props}
            />

            {!inputValue && (
              <div className="absolute inset-0 flex items-center pointer-events-none overflow-hidden">
                <span
                  className="text-zinc-500 text-[15px] transition-all duration-300 ease-in-out"
                  style={{
                    animation: isAnimating
                      ? "placeholder-slide-up 0.3s ease-in-out forwards"
                      : "placeholder-slide-in 0.3s ease-in-out forwards",
                  }}
                >
                  {placeholders[currentIndex]}
                </span>
              </div>
            )}
          </div>

          <div className="flex items-center gap-2 pr-3">
            {inputValue ? (
              <button
                type="button"
                onClick={handleClear}
                className="flex items-center justify-center w-6 h-6 text-zinc-500 hover:text-zinc-300 transition-colors cursor-pointer"
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="16"
                  height="16"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                >
                  <path d="M18 6 6 18" />
                  <path d="m6 6 12 12" />
                </svg>
              </button>
            ) : (
              <kbd className="hidden sm:flex items-center gap-1 px-2 py-1 text-xs text-zinc-500 bg-zinc-800/50 border border-zinc-700/50 rounded-md">
                <span className="text-[13px]">⌘</span>
                <span>K</span>
              </kbd>
            )}
          </div>
        </div>
      </form>
    </>
  );
}

export function SearchBarMinimal({
  placeholders = ["Search...", "Type something...", "What are you looking for?"],
  interval = 3000,
  onChange,
  onSubmit,
  className = "",
  ...props
}) {
  const [currentIndex, setCurrentIndex] = useState(0);
  const [isAnimating, setIsAnimating] = useState(false);
  const [inputValue, setInputValue] = useState("");
  const [isFocused, setIsFocused] = useState(false);

  useEffect(() => {
    if (inputValue || isFocused) return;

    const timer = setInterval(() => {
      setIsAnimating(true);
      setTimeout(() => {
        setCurrentIndex((prev) => (prev + 1) % placeholders.length);
        setIsAnimating(false);
      }, 300);
    }, interval);

    return () => clearInterval(timer);
  }, [placeholders.length, interval, inputValue, isFocused]);

  const handleChange = (e) => {
    setInputValue(e.target.value);
    onChange?.(e);
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    onSubmit?.(inputValue);
  };

  return (
    <>
      <style>{searchBarStyles}</style>
      <form onSubmit={handleSubmit} className={`relative w-full max-w-md ${className}`}>
        <div
          className={`
            relative flex items-center w-full h-11
            bg-zinc-950
            border border-zinc-800
            rounded-lg
            transition-all duration-200
            hover:border-zinc-700
            focus-within:border-zinc-600
            focus-within:ring-2 focus-within:ring-zinc-800
          `}
        >
          <div className="flex items-center justify-center w-10 text-zinc-600">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="16"
              height="16"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              strokeWidth="2"
              strokeLinecap="round"
              strokeLinejoin="round"
            >
              <circle cx="11" cy="11" r="8" />
              <path d="m21 21-4.3-4.3" />
            </svg>
          </div>

          <div className="relative flex-1 h-full">
            <input
              type="text"
              value={inputValue}
              onChange={handleChange}
              onFocus={() => setIsFocused(true)}
              onBlur={() => setIsFocused(false)}
              className="w-full h-full bg-transparent text-zinc-300 text-sm outline-none placeholder-transparent pr-3"
              {...props}
            />

            {!inputValue && (
              <div className="absolute inset-0 flex items-center pointer-events-none overflow-hidden">
                <span
                  className="text-zinc-600 text-sm transition-all duration-300 ease-in-out"
                  style={{
                    animation: isAnimating
                      ? "placeholder-slide-up 0.3s ease-in-out forwards"
                      : "placeholder-slide-in 0.3s ease-in-out forwards",
                  }}
                >
                  {placeholders[currentIndex]}
                </span>
              </div>
            )}
          </div>
        </div>
      </form>
    </>
  );
}
