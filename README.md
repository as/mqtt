# MQTT
A prototype implementation of mqtt to research how _creative_ protocols can be generated with wire9. I discourage the use of this package as a broker or client at this time, so here are some notes instead of documentation:

# Background
Wire9 (github.com/as/wire9)is a protocol generator created some years ago. It has its benefits:
- Generates arbitrarily-nested conformant types
- Based on a sucessful syntax (see draw(3) from plan9)
- It is good at generating specific protocols
  - Started because of Microsoft COM/DCOM
  - Later used to implement draw(3)

But it has issues
- Riddled with bugs, but it's a custom lexer and parser so...
- Beyond that, there are actual bugs (the hard ones)

# The real issues
- Real protocols are strange
  - Not all data are big endian
  - Not all data are alligned
  - Not all data specifications make sense
- Designed by creative people
- Unbounded creativity fosters complexity
  - Varints 
  - Existential values
  - Fan-Out tables
  - Bit Aligned values
  - Variable length headers
  - Computed Conformant Types
  - Musts, Mays, Shoulds
  
# Varints
 - Cut a number into 7 bit pieces
 - Put a cut into an 8 bit byte
 - If theres a cut left, mark the 8th bit '1'
 - Advantage: Saves space when the standard deviation for all numeric messages is high
 - Disadvantage: Time/memory trade off
 - Disadvantage: Implementation complexity
 
 Go did not have a reader/writer type for varints, so I made one. Although there is
 an encoding/binary package that allows slices to become varints, lack of a streaming
 API complicated wire9 at first.
 
 # Existential Values
  - Things like 
    - "if this flag is this value, then this data will exist, otherwise, this next value exists... but only if this other flag is this value..."
    - MQTT likes to do this
    
 # Fan-Out tables / Jump tables
  It seems every protocol designer has their own name for this. It's a table of offsets for specific data in a message.
  - The idea is that once you have the offset you can sub out the reading/processing to another process/thread/goroutine
  - Actually useful when done right
  - Actually not usefull when
    - The table is the last piece of data in the message
    - The table is null
    - The table has overlapping offsets
    - The table's offsets are too large
    - The API is intended for real time streaming
    
 # Byte Unaligned Types
  - If you have ever created a video encoder, I am so sorry.
  - See the GO standard library's image/jpeg for a very small sample of this
  - The headers are based on a stream of bits
   - There are logical reasons for this in the AV world (we assume)
  
 # Variable length headers
   Headers that require processing before it is known how long they are
   - MQTT
    - "I have a static header."
     - And depending on the message I also have this varint describes the length of another header (not really)
 
 # Computed Offsets
   An exaggerated example of what this is

   - Find the length of this conformant type
   - Subtract it from the length of this other type
   - Multiply it by the cunningham chain in the last two messages [mod] the TCP window size
   - That's what the length of the next message ```MUST``` be
   - Unless its a ```RFC 756``` packet
   - Then it ```MAY``` be ```128```

   Actual example
   
   - MQTT has two header types
    - Static: Never changes length
    - Dynamic: Size varies, and the length that describes it is a varint. The varint also includes the size of a potential payload.
    - This payload conforms to a value that is computed
    
# Computed Conformant Types
  Normal conformant types are normally described by a type whose value (normally) preceeds the conformant type
  - But with computed conformant types
    - The length isn't concrete and needs to be computed
    - See computed offsets, except instead of a positional its a number that describes the width of a structure
    - But computation creates a functional dependency
    - Wire9 will not scan a packet twice on read or write; it is not a graph traverser;

# Musts, Mays, Shoulds
>The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT", "SHOULD", "SHOULD NOT", "RECOMMENDED",  "MAY", and "OPTIONAL" in this document are to be interpreted as described in RFC 2119.

The text above is from the commonly-mocked RFC 2119. But why do these phrases confuse?

Each key word can be categorized into three buckets
 - Affirmative: MUST, REQUIRED, SHOULD, SHALL
 - Negative: MUST NOT, SHOULD NOT, SHALL NOT
 - Uncertain: RECOMMENDED, MAY, OPTIONAL
 
Of these, four key words have no inverse: REQUIRED, RECOMMENDED, MAY, OPTIONAL 

 In some key words, the urgency depends on the context. But context sensitivity seems unconducive to a clear spec.

 - You should stay away from those downed power lines
 - You shall stay away from those downed power lines
 - You must stay away from those downed power lines
 
 Often, a negation is more serious
 
 - You should not go near that lion
 - You shall not go near that lion
 - You must not go near that lion
 
 Too many synonyms, too many options. Don't agree with my choice of buckets above? Well that's exactly my point. These words are unformalized.
 
# MQTT
  - A rare protocol
   - Small but complex
   - A perfect case study
   - If wire9 can generate this protocol, it will become better at being a protocol generator
    

   
  
