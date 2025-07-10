# Flappy Bird in Go - System Plan

## First Pass: Key Systems Overview

1. **Game Window System**

   - **Initialization**:
     - Use Ebiten game library for window management
     - Set fixed window dimensions (e.g., 800x600)
     - Initialize game objects and state variables
   - **Rendering Loop**:
     - Implement Ebiten's Game interface (Update, Draw, Layout methods)
     - Target 60 FPS for smooth gameplay
     - Coordinate draw order (background first, then game objects)
   - **Frame Management**:
     - Use delta time for frame-rate independent movement
     - Handle vsync settings
     - Manage resource loading
   - **Inter-system Coordination**:
     - Provide access to input system for keypress detection
     - Pass game state to rendering system
     - Coordinate with game state system for screen transitions
   - **Key Math**:
     - Coordinate system: Origin at top-left corner (0,0)
     - Frame timing calculations: dt = 1/60 ≈ 0.0167s per frame
     - Window dimension constants for collision boundaries

2. **Bird System**

   - **Physics Simulation**:
     - Position (x,y) with initial center-screen placement
     - Velocity (vertical speed) affected by gravity (600px/s²)
     - Jump impulse (-300px/s velocity on flap)
     - Terminal velocity limit (max 400px/s downward)
   - **Movement Logic**:
     - Update position: y += velocity \* dt
     - Apply gravity: velocity += gravity \* dt
     - On flap: set velocity to jump force
     - Constrain vertical movement within screen bounds
   - **Rotation Mechanics**:
     - Base rotation on current velocity
     - Rotate between +30° (max upward) to -90° (max downward)
     - Smooth rotation transitions using lerp (0.1 factor)
   - **Visual Representation**:
     - 40x40px rectangle centered at position
     - Color: yellow (RGB 255, 255, 0)
     - Apply rotation transform before drawing
   - **Inter-system Dependencies**:
     - Input system triggers flap events
     - Collision system requires position/size for detection
     - Game state system pauses updates when game over
     - Rendering system needs transformed draw coordinates
   - **Key Math**:
     - Velocity integration: y = y₀ + v*t + ½a*t²
     - Rotation angle: angle = clamp(velocity \* 0.15, -90, 30)
     - Bounding box: x±20, y±20 for collision checks

3. **Pipe System**

   - **Pipe Generation**:
     - Spawn new pipe pairs every 3 seconds
     - Fixed horizontal spacing (300px between pairs)
     - Gap size of 150px vertical space
     - Initial spawn position at right edge (x=800)
   - **Movement Mechanics**:
     - Constant leftward movement at 200px/s
     - Update position: x -= pipeSpeed \* dt
     - Sync top/bottom pipe movement
   - **Pipe Recycling**:
     - Remove pipes when x < -100 (off-screen left)
     - Reuse pipe objects by resetting position
     - Maintain 3 active pipe pairs simultaneously
   - **Randomization**:
     - Generate random vertical gap positions
     - Minimum 100px margin from top/bottom edges
     - Smooth height transitions between pairs
   - **Visual Representation**:
     - Top pipe: 80px wide, height from top to gap
     - Bottom pipe: 80px wide, height from gap to bottom
     - Green color (RGB 0, 255, 0)
   - **Inter-system Dependencies**:
     - Collision system requires pipe geometries
     - Score system tracks passed pipes
     - Game state system pauses movement when game over
     - Rendering system draws pipes behind bird
   - **Key Math**:
     - Gap center Y = rand(100, windowHeight-100-gapSize)
     - Top pipe height = gapY - (gapSize/2)
     - Bottom pipe y = gapY + (gapSize/2)
     - Pipe spacing = 300px / (pipeSpeed \* spawnInterval)

4. **Collision System**

   - **Bird-Pipe Collision**:
     - Axis-Aligned Bounding Box (AABB) detection
     - Bird bounds: position ±20px (40x40 square)
     - Pipe bounds: top pipe rectangle, bottom pipe rectangle
     - Check overlap between bird and any pipe rectangle
     - Consider pipe width and current position
   - **Ground Collision**:
     - Detect if bird's bottom edge (y+20) >= ground level (windowHeight)
     - Ground level at bottom of screen (y=600)
   - **Ceiling Collision**:
     - Detect if bird's top edge (y-20) <= 0
     - Instant game over on ceiling hit
   - **Gap Passing Detection**:
     - Trigger when bird center passes pipe center (x position)
     - Track which pipes have been scored to avoid duplicates
     - Use pipe ID to prevent multiple scoring
   - **Collision Response**:
     - Set game state to "Game Over" on any collision
     - Freeze game objects
     - Trigger game over screen
   - **Inter-system Dependencies**:
     - Requires bird position/size from Bird System
     - Requires pipe positions/sizes from Pipe System
     - Notifies Game State System on collision
     - Triggers Score System on successful gap pass
   - **Key Math**:
     - AABB collision: birdRight > pipeLeft && birdLeft < pipeRight && birdBottom > pipeTop && birdTop < pipeBottom
     - Gap scoring: if birdX > pipeCenterX && !pipe.scored
     - Ground collision: if birdY + 20 >= groundY

5. **Score System**

   - **Scoring Logic**:
     - Increment by 1 when bird passes through pipe gap
     - Only count each pipe pair once
     - Triggered by Collision System's gap detection
   - **Display Mechanism**:
     - Render current score top-center of screen
     - Use large, readable font (e.g., 48pt sans-serif)
     - White color for contrast against background
   - **High Score Management**:
     - Store high score in file (highscore.txt)
     - Load high score on game start
     - Update when current score exceeds it
     - Display high score on game over screen
   - **Game Over Display**:
     - Show "Game Over" message
     - Display final score and high score
     - Provide restart instructions
   - **Inter-system Dependencies**:
     - Receives scoring events from Collision System
     - Renders through Game Window System
     - Game State System controls when to display scores
     - Input System handles restart triggers
   - **Key Math**:
     - Position: x = (windowWidth - textWidth)/2, y = 80
     - File operations: read/write highscore.txt
     - Formatting: "Score: {current} | Best: {high}"

6. **Game State System**

   - **States**:
     - Start: Display title and instructions
     - Play: Active gameplay
     - GameOver: Show score and restart prompt
   - **State Variables**:
     - Current state enum (start/play/gameOver)
     - Game objects visibility flags
     - Control enabling flags
   - **Transitions**:
     - Start → Play: On spacebar press
     - Play → GameOver: On collision detection
     - GameOver → Start: On spacebar press (restart)
   - **State-Specific Logic**:
     - Start: Only background and title visible
     - Play: All game objects active, physics running
     - GameOver: Freeze objects, display score screen
   - **Initialization**:
     - Start in "Start" state
     - Reset all game objects on restart
   - **Inter-system Dependencies**:
     - Input System triggers transitions
     - Collision System triggers GameOver
     - Controls rendering for each state
     - Coordinates Score System display
   - **Key Considerations**:
     - State-specific update logic
     - State-specific rendering
     - Reset mechanics between games

7. **Input System**
   - **Key Press Handling**:
     - Detect spacebar presses for bird flap
     - Single press per flap (no auto-repeat)
     - Also handle restart input (spacebar in GameOver state)
   - **Touch Input**:
     - Detect touch/click events for mobile support
     - Map to same flap/restart functionality
   - **Restart Mechanism**:
     - Clear game state on restart
     - Reset all game objects
     - Return to Start state
   - **Input Buffering**:
     - Allow input even during state transitions
     - Prevent input during freeze states
   - **Inter-system Dependencies**:
     - Triggers Bird System flap events
     - Controls Game State transitions
     - Coordinates with Score System for restart
   - **Key Considerations**:
     - Platform-agnostic input handling
     - Input cooldown to prevent multiple flaps
     - Visual feedback for user inputs
